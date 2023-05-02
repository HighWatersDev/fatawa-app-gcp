import React, {useState, useEffect} from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import {Routes, Route} from 'react-router-dom';
import HomePage from './page/HomePage';
import LoginPage from './page/LoginPage';
import DocumentsPage from './page/DocumentsPage';

function App() {
  return (
      <Router>
        <div className="App bg-primary">
          <section>
            <div>
              <Routes>
                <Route path="/home" element={< HomePage />}/>
                <Route path="/fatawa" element={< DocumentsPage />}/>
                <Route path="/login" element={<LoginPage/>}/>
              </Routes>
            </div>
          </section>

        </div>
      </Router>
  );
}

export default App;