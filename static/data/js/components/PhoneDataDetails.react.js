var React = require('react');

var PhoneDataDetails = React.createClass({
    render: function() {
        if (!this.props.phoneDataItem) {
            return (
                <span>No entry selected</span>
            );
        }

        var phoneData = this.props.phoneDataItem;
        var contactNodes = phoneData.Contacts.map(function(contact) {
            return (
                <tr>
                    <td>{contact.DisplayName}</td>
                    <td>{contact.Number}</td>
                </tr>
            );
        });

        return (
            <table>
                <thead>
                    <tr>
                        <th>DisplayName</th>
                        <th>Number</th>
                    </tr>
                </thead>
                <tbody>
                {contactNodes}
                </tbody>
            </table>
        );
    }
});

module.exports = PhoneDataDetails;
