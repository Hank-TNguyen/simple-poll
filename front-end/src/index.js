import React from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import './index.css';
import App from './App';
import ViewPoll from './viewpoll';
// Remove the below line if you don't have reportWebVitals.js
import reportWebVitals from './reportWebVitals';


const container = document.getElementById('root');
const root = createRoot(container);

root.render(
  <React.StrictMode>
    <BrowserRouter>
      <Routes>
        {/* Home / Default route */}
        <Route path="/" element={<App />} />

        {/* Route to view a single poll by ID */}
        <Route path="/poll/:pollId" element={<ViewPoll />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);

// If removed, also remove or comment out the reportWebVitals() call.
reportWebVitals(console.log);