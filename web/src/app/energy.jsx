var React = require("react/addons");
var mui = require("material-ui");
var ThemeManager = mui.Styles.ThemeManager();
ThemeManager.setTheme(ThemeManager.types.DARK);
var Paper = mui.Paper,
    AppBar = mui.AppBar,
    IconButton = mui.IconButton,
    FontIcon = mui.FontIcon;

var BackButton = require("components/back_button");

var Energy = React.createClass({
  childContextTypes: {
    muiTheme: React.PropTypes.object
  },

  getChildContext: function() {
    return {
      muiTheme: ThemeManager.getCurrentTheme()
    };
  },

  render: function() {
    return (
      <div>
        <AppBar title="fingerpoken" iconElementLeft={<BackButton/>} />
        <div className="page">
          <Paper> Hello </Paper>
        </div>
      </div>
    );
  }

});

module.exports = Energy;
