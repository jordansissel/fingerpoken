var React = require("react/addons");
var mui = require("material-ui");
var ThemeManager = mui.Styles.ThemeManager();
ThemeManager.setTheme(ThemeManager.types.DARK);
var Paper = mui.Paper,
    AppBar = mui.AppBar,
    TextField = mui.TextField,
    IconButton = mui.IconButton,
    FlatButton = mui.FlatButton,
    FontIcon = mui.FontIcon;
var zws = require("lib/zws");

var BackButton = require("components/back_button");

var Pushover = React.createClass({
  getInitialState: function() {
    return {
      mdc: new zws.MajordomoClient("ws://" + document.location.hostname + ":8111/zws/1.0"),
      messageErrorText: "required"
    };
  },
  childContextTypes: {
    muiTheme: React.PropTypes.object
  },

  getChildContext: function() {
    return {
      muiTheme: ThemeManager.getCurrentTheme()
    };
  },

  messageChange: function(e) {
    if (e.target.value !== "" && this.state.messageErrorText !== "") {
      this.setState({ messageErrorText: "" })
    } else if (e.target.value == "" && this.state.messageErrorText == "") {
      this.setState({ messageErrorText: "required" })
    }
  },

  send: function(e) {
    m = {
      message: this.refs.message.getValue()
    }

    var title = this.refs.title.getValue();
    if (title !== "") {
      m.title = title;
    }

    var rpc = {
      method: "Pushover.Send",
      params: [ m ],
      id: { id: 1 /* TODO(sissel): generate rpc id */, ts: Date.now() },
    }
    var start = Date.now();
    console.log(rpc);
    this.state.mdc.send("pushover", JSON.stringify(rpc), function(service, response) { 
      console.log("RPC latency: " + (Date.now() - start) + "ms");
      console.log(response) 
    });
  },

  render: function() {
    return (
      <div>
        <AppBar title="fingerpoken - pushover" iconElementLeft={<BackButton/>} />
        <div>
          <div> <TextField ref="title" hintText="Title (optional)" floatingLabelText="Title" /> </div>
          <div> <TextField ref="message" hintText="Message" floatingLabelText="Message" onChange={this.messageChange} multiLine={true} errorText={this.state.messageErrorText}/> </div>
          <div> <FlatButton label="Send" onClick={this.send}/> </div>
        </div>
      </div>
    );
  }

});

module.exports = Pushover;
