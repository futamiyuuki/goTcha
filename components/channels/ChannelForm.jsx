import React, {Component} from 'react';

class ChannelForm extends Component {
  submitChannel(e) {
    e.preventDefault();
    if (!this.props.connected) return;
    this.props.addChannel(this.refs.channel.value);
    this.refs.channel.value = '';
  }

  render() {
    return (
      <form onSubmit={this.submitChannel.bind(this)}>
        <div className="form-group">
          <input 
            className="form-control" 
            placeholder="Add New Channel #Ex. HALP"
            type="text" 
            ref="channel" />
        </div>
      </form>
    );
  }
}

export default ChannelForm;
