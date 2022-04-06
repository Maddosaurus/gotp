import './App.css';
import React from 'react'

function App() {
  return (
    <div className="App">
      <Home />
      <OTPEntryDetail />
    </div>
  );
}

class Home extends React.Component {
  render() {
    return (
      <>
        <h1>Pallas</h1>
        <p>The self-hosted OTP sync suite</p>
      </>
    );
  }
  componentDidMount() {
    fetch("http://localhost:8081/v1/otp/entries")
    .then (res => res.json())
    .then((data) => {
      this.setState({entries: data})
    })
    .catch(console.log)
  }
}

function OTPEntryDetail() {
  return (
    <>
      <h1>Entry Title</h1>
      <small>Updated At: 1970-01-01T13:05:25</small>
      <h3>123456</h3>
    </>
  )
}

export default App;
