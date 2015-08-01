var Alert = ReactBootstrap.Alert;
var Button = ReactBootstrap.Button;
var Input = ReactBootstrap.Input;

// Class to display a single mapping
var ShortItem = React.createClass({
  render: function() {
    i = this.props.item
    modified = i.modified_date ? moment(i.modified_date).format("YYYY/M/D, h:mm:ss") : ""
    expire = i.expire_date ? moment(i.expire_date).format("YYYY/M/D, h:mm:ss") : ""
    return (
      <div className="short-item">
        <span className="slug">{i.slug}</span>
        <span className="long-url">{i.long_url}</span>
        <span className="modified-date">{modified}</span>
        <span className="expire-date">{expire}</span>
      </div>
    );
  }
});

// Overall control class
var Shortener = React.createClass({
  shortList: [],

  shortenURL: function(e) {
    e.preventDefault()
    console.log("this ", this.refs.slugToAdd.getInputDOMNode().value)
    $.ajax({
      url: "/shorten",
      dataType: "json",
      data: {
        slug: this.refs.slugToAdd.getInputDOMNode().value,
        long_url: this.refs.longURLToAdd.getInputDOMNode().value
      },
      type: "POST",
      success: function(data) {
        flash = <Alert bsStyle="success">Successfully linked URL</Alert>
        console.log("successful post, flash: " , flash);
        this.setState({flash: flash});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error("error linking URL", xhr.responseJSON);
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
        console.log("fl2: ", this.state.flash);
        this.setState({flash: this.state.flash, shortList: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error("error getting list", xhr.responseJSON);
        this.setState({flash: <Alert bsStyle="danger">Error loading data: {xhr.responseJSON.error}</Alert>});
      }.bind(this)
    });
  },

  getInitialState: function() {
    this.refreshList();
    return {shortList: []};
  },

  render: function() {
    existingItems = [];
    _.each(this.state.shortList, function(shortItem) {
      console.log("item: ", shortItem)
      si = (<li><ShortItem item={shortItem} /></li>)
      existingItems.push(si)
    });
    return (
      <div>
        {this.state.flash}
        <form className="add-url" >
          <Input ref="slugToAdd" className="slug-to-add" type="text" defaultValue="short"></Input>
          <Input ref="longURLToAdd" className="long-url-to-add" type="text" defaultValue="http://example.com/lonnnnnnnnnnnng"></Input>
          <Button onClick={this.shortenURL}>Shorten!</Button>
        </form>
        <ul>
          {existingItems}
        </ul>
      </div>
    );
  }
});

React.render(<Shortener />, $("#shortener")[0]);
