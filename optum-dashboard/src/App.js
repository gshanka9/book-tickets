import React, { useEffect, useState } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [errors, setErrors] = useState([]);

  useEffect(() => {
      const baseURL = 'http://localhost:8080/errors'
    axios.get(baseURL)
        .then(response => {
            console.log("response from server: ", response)
          if (response.status === 200) {
            setErrors(response.data);
          }
          else{
              console.log("No errors returned from the server")
          }
        })
        .catch(error => console.error('Error fetching log entries:', error));
  }, []);

  return (
      <div className="App">
        <Navbar />
        <div className="container">
          <h1>Error Log Dashboard</h1>
          <ErrorGrid errors={errors} />
        </div>
      </div>
  );
}

function Navbar() {
  return (
      <nav className="navbar">
        <span className="navbar-brand">Optum</span>
      </nav>
  );
}

function ErrorGrid({ errors }) {
  if (errors.length === 0) {
    return <p>No errors found.</p>;
  }

  return (
      <table className="error-grid">
        <thead>
        <tr>
          <th>Error Message</th>
          <th>File</th>
          <th>Line Number</th>
          <th>Author</th>
        </tr>
        </thead>
        <tbody>
        {errors.map((error, index) => (
            <tr key={index}>
              <td>{error.msg}</td>
              <td>{error.file}</td>
              <td>{error.line}</td>
              <td>{error.author}</td>
            </tr>
        ))}
        </tbody>
      </table>
  );
}

export default App;
