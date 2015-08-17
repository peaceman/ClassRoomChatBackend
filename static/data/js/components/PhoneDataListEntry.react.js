var React = require('react');
var PhoneDataActions = require('../actions/PhoneDataActions');
var classNames = require('classnames');

var PhoneDataListEntry = React.createClass({
    propTypes: {
        phoneData: React.PropTypes.object.isRequired
    },
    render: function() {
        var phone = this.props.phoneData;

        var tableRowClasses = classNames({
            active: phone.isSelected
        });

        var emailNodes = phone.Emails.map(function(email) {
            return (
                <span className="email-list-entry">{email}</span>
            );
        });

        return (
            <tr onClick={this._onToggleSelection} className={tableRowClasses}>
                <td>{emailNodes}</td>
                <td>{phone.PhoneNumber}</td>
                <td>{phone.Build.Manufacturer} {phone.Build.Model}</td>
                <td>
                    <button className="btn btn-danger btn-xs" onClick={this._onDestroyClick}>Delete</button>
                </td>
            </tr>
        );
    },
    _onToggleSelection: function() {
        PhoneDataActions.toggleSelection(this.props.phoneData.id);
    },
    _onDestroyClick: function() {
        PhoneDataActions.destroy(this.props.phoneData.id);
        return false; // end further onClick event handling
    }
});

module.exports = PhoneDataListEntry;
