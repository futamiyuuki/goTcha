import React, {Component} from 'react';
import Channel from './Channel.jsx';

class ChannelList extends Component {
  render() {
    const channels = this.props.channels.map((channel) => 
      <Channel 
        channel={channel} 
        {...this.props}
        key={channel.id} />);

    return (
      <ul>{channels}</ul>
    );
  }
}

export default ChannelList;
