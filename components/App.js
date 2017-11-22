import React, { Component } from 'react';
import MessageSection from './messages/MessageSection.jsx';
import ChannelSection from './channels/ChannelSection.jsx';
import UserSection from './users/UserSection.jsx';
import Socket from '../util/socket.js';

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      messages: [],
      channels: [],
      users: [],
      currentChannel: null,
      connected: false,
    }
  }

  componentDidMount() {
    console.log(process.env.WS_PORT);
    const socket = this.socket = new Socket(new WebSocket(process.env.WS_PORT || 'wss://young-garden-89860.herokuapp.com/ws'));
    socket.on('connect', this.onConnect.bind(this));
    socket.on('disconnect', this.onDisconnect.bind(this));
    socket.on('message add', this.onAddMessage.bind(this));
    socket.on('channel add', this.onAddChannel.bind(this));
    socket.on('user add', this.onAddUser.bind(this));
    socket.on('user edit', this.onEditUser.bind(this));
    socket.on('user remove', this.onRemoveUser.bind(this));
  }

  onConnect() {
    console.log('Connected to server!');
    this.setState({connected: true});
    this.socket.send('channel subscribe');
    this.socket.send('user subscribe');
  }

  onDisconnect() {
    this.setState({connected: false});
  }

  onAddMessage(message) {
    console.log('on add message');
    const {messages} = this.state;
    messages.push(message);
    this.setState({messages});
  }

  onAddChannel(channel) {
    console.log('channel added:', channel);
    const {channels} = this.state;
    channels.push(channel);
    this.setState({channels});
  }

  onAddUser(user) {
    console.log('user added', user);
    const {users} = this.state;
    users.push(user);
    this.setState({users});
  }

  onEditUser(target) {
    console.log(target);
    let {users} = this.state;
    users = users.map(user => target.id === user.id ? target : user);
    this.setState({users});
  }

  onRemoveUser(target) {
    let {users} = this.state;
    users = users.filter(user => target.id !== user.id);
    this.setState({users});
  }

  addMessage(content) {
    const {currentChannel} = this.state;
    this.socket.send('message add', {channelId: currentChannel.id, content});
  }

  addChannel(name) {
    const {channels} = this.state;
    this.socket.send('channel add', {name});
  }

  setChannel(currentChannel) {
    this.setState({currentChannel});
    this.socket.send('message unsubscribe');
    this.setState({messages: []});
    this.socket.send('message subscribe', {channelId: currentChannel.id});
  }

  setUser(currentUserName) {
    this.socket.send('user edit', {currentUserName});
  }

  render() {
    return (
      <div className="app">
        <MessageSection
            {...this.state}
            addMessage={this.addMessage.bind(this)} />
        <div className="nav">
          <ChannelSection
            {...this.state}
            addChannel={this.addChannel.bind(this)}
            setChannel={this.setChannel.bind(this)} />
          <UserSection
            {...this.state}
            setUser={this.setUser.bind(this)} />
        </div>
      </div>
    );
  }
}

export default App;
