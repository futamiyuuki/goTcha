import React, {Component} from 'react';
import MessageList from './MessageList.jsx';
import MessageForm from './MessageForm.jsx';

class MessageSection extends Component{
  render() {
    const {currentChannel} = this.props;
    return (
      <div className='messages-container panel panel-default'>
        <div className='panel-heading'>{currentChannel && currentChannel.name || 'Select A Channel'}</div>
        <div className='panel-body messages'>
          <MessageList {...this.props} />
          <MessageForm {...this.props} />
        </div>
      </div>
    )
  }
}

export default MessageSection
