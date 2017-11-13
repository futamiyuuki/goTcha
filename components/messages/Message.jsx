import React from 'react';

const Message = ({message}) => {
  return (
    <li className='message'>
      <div className='author'>
        <strong>{message.author}</strong> 
        <i className='timestamp'>{message.createdAt}</i>
      </div>
      <div className='body'>{message.content}</div>
    </li>
  );
};

export default Message;
