import React, {Component} from 'react';

class UserForm extends Component {
  submitUser(e) {
    e.preventDefault();
    this.props.setUser(this.refs.user.value);
    this.refs.user.value = '';
  }

  render() {
    return (
      <form onSubmit={this.submitUser.bind(this)}>
        <div className="form-group">
          <input 
            className="form-control" 
            placeholder="Set User Name #Ex. TreeFiddy1337"
            type="text" 
            ref="user" />
        </div>
      </form>
    );
  }
}

export default UserForm;

