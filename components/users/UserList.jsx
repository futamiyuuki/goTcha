import React, {Component} from 'react';
import User from './User.jsx';

class UserList extends Component {
  render() {
    const users = this.props.users.map((user) => 
      <User 
        user={user}
        key={user.id} />);

    return (
      <ul>{users}</ul>
    );
  }
}

export default UserList;
