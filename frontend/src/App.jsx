import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Navbar from "./components/Navbar";
import Home from "./components/HomeComponents/Home";
import Login from "./components/Login";
import Notes from "./components/NotesComponents/Notes";
import { ToastProvider } from "react-toast-notifications";

// import Footer from './components/Footer';

const App = () => (
  <>
    <Router>
      <Navbar />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/notes" element={<Notes />} />
        <Route
          path="/login"
          element={
            <ToastProvider>
              <Login />
            </ToastProvider>
          }
        />
      </Routes>
    </Router>
    {/* <Footer /> */}
  </>
);

export default App;
