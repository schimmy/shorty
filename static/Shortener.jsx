var Alert = ReactBootstrap.Alert;
var Button = ReactBootstrap.Button;
var Input = ReactBootstrap.Input;

// Class to display a single mapping
var ShortItem = React.createClass({
  deleteURL: function(e) {
    e.preventDefault();
    this.props.deleteURL(this.props.item.slug);
  },

  render: function() {
    i = this.props.item
    //expire = i.expire_date ? moment(i.expire_date).format("YYYY/M/D") : ""
    shortURL = "http://go/" + i.slug
    return (
      <tr className="short-item">
        <td className="slug"><a href={shortURL}>go/{i.slug}</a></td>
        <td className="long-url"><a href={i.long_url}>{i.long_url}</a></td>
        <td className="owner">{i.owner}</td>
        <td><Button onClick={this.deleteURL} className="glyphicon glyphicon-remove"></Button></td>
      </tr>
    );
  }
});

// Overall control class
var Shortener = React.createClass({
  shortList: [],

  shortenURL: function(e) {
    e.preventDefault()
    slug = this.refs.slugToAdd.getInputDOMNode().value
    owner = this.refs.ownerToAdd.getInputDOMNode().value
    $.ajax({
      url: "/shorten",
      dataType: "json",
      data: {
        slug: slug,
        long_url: this.refs.longURLToAdd.getInputDOMNode().value,
        owner: owner
      },
      type: "POST",
      success: function() {
        newLink = "http://go/"+ slug;
        this.refreshList();
        flash = <Alert bsStyle="success">Successfully linked URL: <a href={newLink}>{newLink}</a></Alert>
        this.setState({currentOwner: owner, flash: flash});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error("error linking URL", xhr.responseJSON);
        this.setState({currentOwner: owner, flash: <Alert bsStyle="danger">{xhr.responseJSON.error}</Alert>});
      }.bind(this)
    });
  },

  deleteURL: function(slug) {
    $.ajax({
      url: "/delete",
      dataType: "json",
      data: { slug: slug },
      type: "POST",
      success: function(data) {
        flash = <Alert bsStyle="success">Successfully deleted URL</Alert>
        console.log("successful post, flash: " , flash);
        this.setState({currentOwner: this.state.currentOwner, flash: flash});
        this.refreshList();
      }.bind(this),
      error: function(xhr, status, err) {
        console.error("error deleting URL", xhr.responseJSON);
        this.setState({currentOwner: this.state.currentOwner, flash: <Alert bsStyle="danger">{xhr.responseJSON.error}</Alert>});
      }.bind(this)
    });
  },

  refreshList: function() {
    $.ajax({
      url: "/list",
      dataType: "json",
      type: "GET",
      success: function(data) {
        this.setState({currentOwner: this.state.currentOwner, flash: this.state.flash, shortList: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error("error getting list", xhr.responseJSON);
        this.setState({currentOwner: this.state.currentOwner, flash: <Alert bsStyle="danger">Error loading data: {xhr.responseJSON.error}</Alert>});
      }.bind(this)
    });
  },

  getInitialState: function() {
    this.refreshList();
    return {shortList: [], currentOwner: "None"};
  },

  render: function() {
    existingItems = [];
    _this = this
    _.each(this.state.shortList, function(shortItem) {
      si = (<ShortItem item={shortItem} deleteURL={_this.deleteURL} />)
      existingItems.push(si)
    });
    return (
      <div>
        {this.state.flash}
        <div className="panel panel-default">
          <div className="panel-heading">Add or Modify a URL</div>
          <form className="add-url" >
            <div className="pre-text">go/</div>
            <Input ref="slugToAdd" className="slug-to-add add-item" type="text" defaultValue="short" required></Input>
            <div className="pre-text">→</div>
            <Input ref="longURLToAdd" className="long-url-to-add add-item" pattern="http.*" type="text" defaultValue="http://example.com/lonnnnnnnnnnnng" required></Input>
            <div className="pre-text">owned by</div>
            <Input ref="ownerToAdd" className="owner-to-add add-item" type="text" defaultValue={this.state.currentOwner} required></Input>
            <Button onClick={this.shortenURL}>Shorten!</Button>
          </form>
        </div>
        <div className="panel panel-default">
          <div className="table-list panel-heading">Existing URLs</div>
            <table className="table">
              <thead><th>Slug</th><th>Long URL</th><th>Owner</th><th></th></thead>
              <tbody>
                {existingItems}
              </tbody>
            </table>
        </div>
      </div>
    );
  }
});

React.render(<Shortener />, $("#shortener")[0]);
