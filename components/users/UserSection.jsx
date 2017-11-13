import React, {Component} from 'react';
import UserList from './UserList.jsx';
import UserForm from './UserForm.jsx';

class UserSection extends Component {
  render() {
    return (
      <div className="support panel panel-primary">
        <div className="panel-heading">Users</div>
        <div className="panel-body users">
          <UserList {...this.props} />
          <UserForm {...this.props} />
        </div>
      </div>
    );
  }
}

export default UserSection;
