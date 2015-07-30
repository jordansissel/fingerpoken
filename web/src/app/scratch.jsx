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

var Scratch = React.createClass({
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
  },

  componentWillUnmount: function() { 
  },

  photoReady: function(e) {
    var upload = e.target.files[0]
    alert("Size " + upload.size);
    if (upload.size > 1 << 20) {
      alert("Size " + upload.size + " is too large");
      console.log("Image is too large")
      e.target.value = undefined;
      e.preventDefault();
    }
  },

  render: function() {
    return (
      <div>
        <AppBar title="Screen Brightness"  iconElementLeft={<BackButton/>} />
        
        <input type="file" accept="image/*" capture onChange={this.photoReady} name="image"/>
      </div>
    );
  }

});

module.exports = Scratch;
