import React from 'react';
import { Link } from "react-router-dom";

class Home extends React.Component {
    constructor(props) {
      super(props);
      this.state = {entries: []};
    }

    entryList(props) {
      const entries = props.entries;
      const listItems = entries.map((entry) =>
        <li key={entry.uuid}><Link key={entry.uuid} to={`/otp/${entry.uuid}`}>{entry.name}</Link></li>
      );
      return (
        <ul>
          {listItems}
        </ul>
      )
    }

    render() {
      return (
        <>
          <h1>Pallas</h1>
          <p>The self-hosted OTP sync suite</p>
          <p>Currently holding {this.state.entries.length} entries:</p>
          <this.entryList entries={this.state.entries}/>
          <p>ToDo: Investigate https://react-redux.js.org/</p>
        </>
      );
    }

    componentDidMount() {
      fetch("http://localhost:8081/v1/otp/entries")
      .then (res => res.json())
      .then((data) => {
        this.setState({entries: data.entries})
      })
      .catch(console.log)
    }
  }


export default Home
