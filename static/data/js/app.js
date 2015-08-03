var React = require('react');

var ClassRoomChatDataApp = require('./components/ClassRoomChatDataApp.react');
var PhoneDataAPI = require('./utils/PhoneDataAPI');

PhoneDataAPI.init('ws://' + location.host + '/data');

React.render(
    <ClassRoomChatDataApp />,
    document.getElementById('classroomchatdataapp')
);
