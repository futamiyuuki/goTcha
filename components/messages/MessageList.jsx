import React, {Component} from 'react';
import ReactDOM from 'react-dom';
import Message from './Message.jsx';

class MessageList extends Component{
  componentDidUpdate() {
    this.scrollToBottom();
  }

  scrollToBottom() {
    const { messageList } = this.refs;
    const scrollHeight = messageList.scrollHeight;
    const height = messageList.clientHeight;
    const maxScrollTop = scrollHeight - height;
    ReactDOM.findDOMNode(messageList).scrollTop = maxScrollTop > 0 ? maxScrollTop : 0;
  }

  render() {
    const messages = this.props.messages.map(message => <Message message={message} key={message.id} />);
    return (
      <ul ref="messageList">{messages}</ul>
    )
  }
}

export default MessageList;