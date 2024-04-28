import axios from "axios";

const API_URL = "http://localhost:3000/";

const register = (username, Passwd, first_name, last_name, email, phone, dob, gender) => {
  return axios.post(API_URL + "register", {
    username,
    Passwd,
    first_name,
    last_name,
    email,
    phone,
    dob,
    gender: parseInt(gender),
    
  });
};


const login = (Username, Passwd) => {
  return axios
    .post(API_URL + "login", {
      Username,
      Passwd,
    })
    .then((response) => {
      if (response.data.token) {
        localStorage.setItem("user", JSON.stringify(response.data));
      }

      return response.data;
    });
};

const logout = () => {
  localStorage.removeItem("user");
  localStorage.removeItem("token");
};

const getCurrentUser = () => {
  const user = JSON.parse(localStorage.getItem("user"));
  const token = localStorage.getItem("token");
  if (user && token) {
    return { ...user, token };
  } else {
    return null;
  }
};
const AuthService = {
  register,
  login,
  logout,
  getCurrentUser,
};

export default AuthService;
