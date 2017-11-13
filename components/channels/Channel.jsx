import React from 'react';

const Channel = ({setChannel, channel, currentChannel}) => {
  const channelStyle = currentChannel === channel ? 'selected inList' : 'inList';
  return (
    <li className={channelStyle} onClick={() => setChannel(channel)} >{channel.name}</li>
  );
};

export default Channel;
