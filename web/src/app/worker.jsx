var React = require("react/addons");
var mui = require("material-ui");
var ThemeManager = mui.Styles.ThemeManager();
ThemeManager.setTheme(ThemeManager.types.DARK);
var AppBar = mui.AppBar,
    IconButton = mui.IconButton,
    Slider = mui.Slider,
    Paper = mui.Paper,
    FontIcon = mui.FontIcon;
var BackButton = require("components/back_button");
var zws = require("lib/zws");

var Worker = React.createClass({
  getInitialState: function() {
    return { };
  },
  childContextTypes: {
    muiTheme: React.PropTypes.object
  },

  getChildContext: function() {
    return {
      muiTheme: ThemeManager.getCurrentTheme()
    };
  },

  componentWillMount: function() { 
    this.worker = new zws.MajordomoWorker("ws://" + document.location.hostname + ":8111/zws/1.0", "webclient");
    this.worker.rpcHandler = {
      ping: function(body, reply) {
        reply.result = body;
      },
      openurl: function(body, reply) {
        console.log("Openning " + body[0].url);
        document.location.href = body[0].url;
      }
    };
  },

  componentWillUnmount: function() { 
  },

  render: function() {
    return (
      <div>
        <AppBar title="Screen Brightness"  iconElementLeft={<BackButton/>} />
      </div>
    );
  }

});

module.exports = Worker;

