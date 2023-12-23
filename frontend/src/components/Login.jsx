import React, { useState } from "react";
import { useToasts } from "react-toast-notifications";
import "./CSS/Login.css";
import LoginImg from "./Assets/login.svg";
import SignupImg from "./Assets/notetaking.svg";
import { useNavigate } from "react-router-dom";

function Login() {
  const [selectedForm, setSelectedForm] = useState("signup");
  const { addToast } = useToasts();
  const navigate = useNavigate();
  const handleSwitchClick = (e) => {
    setSelectedForm(e.target.textContent.toLowerCase());
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const formData = new FormData(e.target);
    const data = {};
    formData.forEach((value, key) => {
      data[key] = value;
    });

    try {
      console.log(data);
      const response = await fetch(
        `https://neko-notesbackendstorage.onrender.com/${selectedForm}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(data),
        }
      );

      const result = await response.json();

      if (response.ok) {
        // Successful login/signup
        addToast("successfull", { appearance: "success" });
        navigate("/notes");
      } else {
        // Display error using toast
        addToast(result.error, { appearance: "error" });
      }
    } catch (error) {
      console.error("API error:", error);
      addToast("An error occurred. Please try again later.", {
        appearance: "error",
      });
    }
  };

  return (
    <div className="account-wrapper">
      <div
        className={`account-login ${selectedForm === "login" ? "" : "hidden"}`}
      >
        <div className="account-info">
          <h1>Login</h1>
          <form onSubmit={handleSubmit}>
            <input
              type="email"
              name="email"
              id="email"
              placeholder="Enter your Email"
            />
            <input
              type="password"
              name="password"
              id="password"
              placeholder="Enter your Password"
            />
            <button type="submit">Submit</button>
          </form>
          <p>
            Don't have an account ?
            <span className="switch" onClick={handleSwitchClick}>
              Signup
            </span>
          </p>
        </div>

        <div className="account-image">
          <img src={LoginImg} alt="" />
        </div>
      </div>

      <div
        className={`account-signup ${
          selectedForm === "signup" ? "" : "hidden"
        }`}
      >
        <div className="account-image">
          <img src={SignupImg} alt="" />
        </div>

        <div className="account-info">
          <h1>Sign Up</h1>
          <form onSubmit={handleSubmit}>
            <input
              type="text"
              name="firstname"
              id="fname"
              placeholder="Enter your First Name"
            />
            <input
              type="text"
              name="lastname"
              id="lname"
              placeholder="Enter your Last Name"
            />
            <input
              type="email"
              name="username"
              id="email"
              placeholder="Enter your Email"
            />
            <input
              type="password"
              name="password"
              id="password"
              placeholder="Enter your Password"
            />
            <button type="submit">Submit</button>
          </form>
          <p>
            Already have an account ?
            <span className="switch" onClick={handleSwitchClick}>
              Login
            </span>
          </p>
        </div>
      </div>
    </div>
  );
}

export default Login;
