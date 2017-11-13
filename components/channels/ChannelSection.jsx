import React, {Component} from 'react';
import ChannelList from './ChannelList.jsx';
import ChannelForm from './ChannelForm.jsx';

class ChannelSection extends Component {
  render() {
    return (
      <div className="support panel panel-primary">
        <div className="panel-heading">Channels</div>
        <div className="panel-body channels">
          <ChannelList {...this.props} />
          <ChannelForm {...this.props} />
        </div>
      </div>
    );
  }
}

export default ChannelSection;
