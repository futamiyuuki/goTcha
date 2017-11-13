import React, {Component} from 'react';

class MessageForm extends Component{
  sendMessage(e){
    e.preventDefault();
    this.props.addMessage(this.refs.message.value);
    this.refs.message.value = '';
  }
  render() {
    let input;
    if(this.props.currentChannel){
      input = (
        <input 
          className='form-control'
          type='text'
          ref='message'
          placeholder='Messages...' />
      )
    }
    return (
      <form onSubmit={this.sendMessage.bind(this)}>
        <div className='form-group'>
          {input}
        </div>
      </form>
    )
  }
}

export default MessageForm
