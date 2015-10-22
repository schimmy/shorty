var Alert = ReactBootstrap.Alert;
var Button = ReactBootstrap.Button;
var Input = ReactBootstrap.Input;

// Class to display a single mapping
var ShortItem = React.createClass({
  deleteURL: function(e) {
    e.preventDefault();
    this.props.deleteURL(this.props.item.slug, this.props.item.long_url);
  },

  render: function() {
    i = this.props.item
    //expire = i.expire_date ? moment(i.expire_date).format("YYYY/M/D") : ""
    shortURL = this.props.protocol + "://" + this.props.domain + "/" + i.slug
    return (
      <tr className="short-item">
        <td className="slug"><a href={shortURL}>{this.props.domain}/{i.slug}</a></td>
        <td className="long-url"><a href={i.long_url}>{i.long_url}</a></td>
        <td className="owner">{i.owner}</td>
        <td className="delete"><Button onClick={this.deleteURL} className="glyphicon glyphicon-remove"></Button></td>
      </tr>
    );
  }
});

// Overall control class
var Shortener = React.createClass({
  shortList: [],
  protocol: "http",
  domain: "go",

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
        newLink = this.state.protocol + "://" + "/" + slug;
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

  deleteURL: function(slug, longURL) {
    $.ajax({
      url: "/delete",
      dataType: "json",
      data: { slug: slug },
      type: "POST",
      success: function(data) {
        flash = <Alert bsStyle="success">Successfully deleted slug: <em>{slug}</em> pointing to URL: <em>{longURL}</em></Alert>
        console.log("successful post, flash: " , flash);
        this.setState({flash: flash});
        this.refreshList();
      }.bind(this),
      error: function(xhr, status, err) {
        console.error("error deleting URL", xhr.responseJSON);
        this.setState({flash: <Alert bsStyle="danger">{xhr.responseJSON.error}</Alert>});
      }.bind(this)
    });
  },

  refreshList: function() {
    $.ajax({
      url: "/list",
      dataType: "json",
      type: "GET",
      success: function(data) {
        this.setState({shortList: data});
      }.bind(this),
      error: function(xhr, status, err) {
        errMsg = "unknown";
        if (xhr.responseJSON !== undefined) {
          errMsg = xhr.responseJSON.error;
        }
        console.error("error getting list: ", errMsg);
        this.setState({flash: <Alert bsStyle="danger">Error loading data: {errMsg}</Alert>});
      }.bind(this)
    });
  },


  getInitialState: function() {
    $.ajax({
      url: "/meta",
      dataType: "json",
      type: "GET",
      success: function(data) {
        this.setState({protocol: data["protocol"], domain: data["domain"], flash: null});
        // now get the actual list of URLs
        this.refreshList();
      }.bind(this),
      error: function(xhr, status, err) {
        console.error("error getting metadata like protocol and domain", xhr.responseJSON);
        this.setState({currentOwner: this.state.currentOwner, flash: <Alert bsStyle="danger">Error loading metadata: {xhr.responseJSON.error}</Alert>});
      }.bind(this)
    });
    return {shortList: [], currentOwner: "None", flash: <Alert bsStyle="info">Loading...</Alert>}
  },

  render: function() {
    existingItems = [];
    _this = this
    _.each(this.state.shortList, function(shortItem) {
      si = (<ShortItem item={shortItem} deleteURL={_this.deleteURL} protocol={_this.state.protocol} domain={_this.state.domain}/>)
      existingItems.push(si)
    });
    return (
      <div id="inner-shortener">
        {this.state.flash}
        <div className="panel panel-default">
          <div className="panel-heading">Add or Modify a URL</div>
          <form className="add-url" >
            <div className="pre-text">{this.state.domain + "/"}</div>
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
