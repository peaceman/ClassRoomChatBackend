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

        var pictureNodes = phoneData.Pictures.map(function(base64WebPString) {
            var src = "data:image/jpeg;base64," + base64WebPString;
            return (
                <li>
                    <img className="col-xs-6 col-md-6" src={src}/>
                </li>
            );
        });

        return (
            <div>
                <ul className="list-unstyled row">
                    {pictureNodes}
                </ul>
                <table className="table table-condensed table-striped">
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
            </div>
        );
    }
});

module.exports = PhoneDataDetails;
