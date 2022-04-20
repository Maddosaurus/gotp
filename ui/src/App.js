import './App.css';
import React from 'react'
import { Routes, Route } from "react-router-dom";
import Home from './components/Home/Home';
import OTPEntry from './components/OTPEntry/OTPEntry'

function App() {
  return (
    <div className="App">
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/otp/:uuid" element={<OTPEntry />}></Route>
      </Routes>
    </div>
  );
}

export default App;
