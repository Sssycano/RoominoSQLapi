import React, { useState } from "react";
import Form from "react-validation/build/form";
import Input from "react-validation/build/input";
import CheckButton from "react-validation/build/button";
import { isEmail } from "validator";
import AuthService from "../services/auth.service";

const genderOptions = [
  { value: 2, label: "Male" },
  { value: 1, label: "Female" },
  { value: 0, label: "Other" },
  { value: 3, label: "Prefer Not to Say" }
];

const required = (value) => {
  if (!value) {
    return <div className="alert alert-danger">This field is required!</div>;
  }
};

const validEmail = (value) => {
  if (!isEmail(value)) {
    return <div className="alert alert-danger">Invalid email address!</div>;
  }
};

const validLength = (value, min, max) => {
  if (value.length < min || value.length > max) {
    return <div className="alert alert-danger">Length must be between {min} and {max} characters!</div>;
  }
};

const Register = () => {
  const [formData, setFormData] = useState({
    username: "",
    first_name: "",
    last_name: "",
    dob: "",
    gender: "",
    email: "",
    phone: "",
    passwd: "",
  });

  const [successful, setSuccessful] = useState(false);
  const [message, setMessage] = useState("");

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleGenderChange = (e) => {
    setFormData({ ...formData, gender: e.target.value });
  };

  const handleRegister = (e) => {
    e.preventDefault();
    setMessage("");
    setSuccessful(false);

    if (formData.username && formData.passwd) {
      AuthService.register(
        formData.username,
        formData.passwd,
        formData.first_name,
        formData.last_name,
        formData.email,
        formData.phone,
        formData.dob,
        formData.gender 
      ).then(
        (response) => {
          setSuccessful(true);
          setMessage("Registration successful!");
        },
        (error) => {
          const resMessage =
            error.response?.data?.message ||
            error.toString();
          if (error.response?.data?.error === "Userexists") {
            setMessage("User exists");
          } else {
            setMessage(resMessage);
          }
          setSuccessful(false);
        }
      );
    }
  };

  return (
    <div className="col-md-12">
      <div className="card card-container">
        <Form onSubmit={handleRegister}>
          {!successful && (
            <div>
              <div className="form-group">
                <label htmlFor="username">Username</label>
                <Input
                  type="text"
                  className="form-control"
                  name="username"
                  value={formData.username}
                  onChange={handleChange}
                  validations={[required, (value) => validLength(value, 1, 20)]}
                />
              </div>

              <div className="form-group">
                <label htmlFor="first_name">First Name</label>
                <Input
                  type="text"
                  className="form-control"
                  name="first_name"
                  value={formData.first_name}
                  onChange={handleChange}
                  validations={[required, (value) => validLength(value, 1, 20)]}
                />
              </div>

              <div className="form-group">
                <label htmlFor="last_name">Last Name</label>
                <Input
                  type="text"
                  className="form-control"
                  name="last_name"
                  value={formData.last_name}
                  onChange={handleChange}
                  validations={[required, (value) => validLength(value, 1, 20)]}
                />
              </div>

              <div className="form-group">
                <label htmlFor="dob">Date of Birth</label>
                <Input
                  type="date"
                  className="form-control"
                  name="dob"
                  value={formData.dob}
                  onChange={handleChange}
                  validations={[required]}
                />
              </div>

              <div className="form-group">
                <label htmlFor="gender">Gender</label>
                <select
                  className="form-control"
                  name="gender"
                  value={formData.gender}
                  onChange={handleGenderChange}
                  validations={[required]}
                >
                  {genderOptions.map(option => (
                    <option key={option.value} value={option.value}>{option.label}</option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label htmlFor="email">Email</label>
                <Input
                  type="text"
                  className="form-control"
                  name="email"
                  value={formData.email}
                  onChange={handleChange}
                  validations={[required, validEmail]}
                />
              </div>

              <div className="form-group">
                <label htmlFor="phone">Phone</label>
                <Input
                  type="text"
                  className="form-control"
                  name="phone"
                  value={formData.phone}
                  onChange={handleChange}
                  validations={[required, (value) => validLength(value, 1, 20)]}
                />
              </div>

              <div className="form-group">
                <label htmlFor="passwd">Password</label>
                <Input
                  type="password"
                  className="form-control"
                  name="passwd"
                  value={formData.passwd}
                  onChange={handleChange}
                  validations={[required, (value) => validLength(value, 3, 40)]}
                />
              </div>

              <div className="form-group">
                <button className="btn btn-primary">Sign Up</button>
              </div>
            </div>
          )}

          {message && (
            <div className="form-group">
              <div className={successful ? "alert alert-success" : "alert alert-danger"}>
                {message}
              </div>
            </div>
          )}

          <CheckButton style={{ display: "none" }} />
        </Form>
      </div>
    </div>
  );
};

export default Register;
