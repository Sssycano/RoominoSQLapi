import React, { useState } from "react";
import UserService from "../services/user.service";

const RegisterPet = () => {
  const [currentPetName, setCurrentPetName] = useState("");
  const [currentPetType, setCurrentPetType] = useState("");
  const [currentPetSize, setCurrentPetSize] = useState("");
  const [message, setMessage] = useState("");

  const handleRegisterPet = async () => {
    try {
      await UserService.registerpet(currentPetName, currentPetType, currentPetSize);
      setMessage("Pet registration successful!");
    } catch (error) {
        console.error('Error updating pet:', error);
        if (error.response && error.response.data.error === 'DUPLICATE_KEY: pet already exists') {
          window.alert('DUPLICATE_KEY: pet already exists');
        } else {
          setMessage(error.response.data.error);
        }
    }
  };

  return (
    <div className="container">
      <h2>Register Pet</h2>
      <div className="form-group">
        <label htmlFor="currentPetName">Pet Name</label>
        <input
          type="text"
          className="form-control"
          id="currentPetName"
          value={currentPetName}
          onChange={(e) => setCurrentPetName(e.target.value)}
        />
      </div>
      <div className="form-group">
        <label htmlFor="currentPetType">Pet Type</label>
        <input
          type="text"
          className="form-control"
          id="currentPetType"
          value={currentPetType}
          onChange={(e) => setCurrentPetType(e.target.value)}
        />
      </div>
      <div className="form-group">
        <label htmlFor="currentPetSize">Pet Size</label>
        <input
          type="text"
          className="form-control"
          id="currentPetSize"
          value={currentPetSize}
          onChange={(e) => setCurrentPetSize(e.target.value)}
        />
      </div>
      <button className="btn btn-primary" onClick={handleRegisterPet}>Register Pet</button>
      {message && <div className="mt-3">{message}</div>}
    </div>
  );
};

export default RegisterPet;
