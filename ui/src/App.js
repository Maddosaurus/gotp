import './App.css';

function App() {
  return (
    <div className="App">
      <Home />
      <OTPEntryDetail />
    </div>
  );
}

function Home() {
  return (
    <>
      <h1>Pallas</h1>
      <p>The self-hosted OTP sync suite</p>
    </>
  );

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
