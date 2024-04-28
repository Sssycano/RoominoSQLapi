import axios from "axios";
import authHeader from "./auth-header";

const API_URL = "http://localhost:3000/";

const getUserBoard = () => {
  return axios.get(API_URL + "profile", { headers: authHeader() })  .then(response => {
    return response.data;
  })
  .catch(error => {
    throw error;
  });
};

const getpet = () => {
  return axios.get(API_URL + "profile/petupdate", { headers: authHeader() })
  .then(response => {
    return response.data;
  })
  .catch(error => {
    throw error;
  });
};

const updatepet = (current_pet_name, current_pet_type, new_pet_name, new_pet_type, new_pet_size) => {
  return axios.post(API_URL + "profile/petupdate", {current_pet_name, current_pet_type, new_pet_name, new_pet_type, new_pet_size}, { headers: authHeader() });
};

const getinterests = (unit_rent_id) => {
  return axios.get(API_URL + "profile/interests", {params: { unit_rent_id }, headers: authHeader(), })
  .then(response => response.data)
  .catch(error => {
    throw error;
  });
};

const getcomplexunitinfo = (unit_rent_id) => {
  return axios.get(API_URL + "profile/complexunitinfo", {params: { unit_rent_id }, headers: authHeader(), })
  .then(response => response.data)
  .catch(error => {
    throw error;
  });
};


const createinterests = (unit_rent_id, roommate_cnt, move_in_date) => {
  return axios.post(API_URL + "profile/interests", {unit_rent_id, roommate_cnt, move_in_date}, { headers: authHeader() });
};


const registerpet = (current_pet_name, current_pet_type, current_pet_size) => {
  return axios.post(API_URL + "profile/petregister", {current_pet_name, current_pet_type, current_pet_size}, { headers: authHeader() });
};
const searchUnits = (company_name, building_name) => {
  return axios.post(API_URL+"profile/unitinfo", { company_name, building_name }, { headers: authHeader() })
    .then(response => {
      return response.data;
    })
    .catch(error => {
      throw error;
    });
};

const searchInterests = (unit_rent_id, roommate_cnt,move_in_date) => {
  return axios.post(API_URL+"profile/searchinterests", { unit_rent_id, roommate_cnt,move_in_date }, { headers: authHeader() })
    .then(response => {
      return response.data;
    })
    .catch(error => {
      throw error;
    });
};

const getinterestsProfile = (user_name) => {
  return axios.get(API_URL + "profile/interestprofile", {params: { user_name }, headers: authHeader(), })
  .then(response => response.data)
  .catch(error => {
    throw error;
  });
};

const avgPrice = (addr_zip_code, bedroom_num,bathroom_num) => {
  return axios.post(API_URL+"profile/avgprice", { addr_zip_code, bedroom_num,bathroom_num }, { headers: authHeader() })
    .then(response => {
      return response.data;
    })
    .catch(error => {
      throw error;
    });
};
const UserService = {
  updatepet,
  getUserBoard,
  searchUnits,
  getpet,
  registerpet,
  getinterests,
  createinterests,
  getcomplexunitinfo,
  searchInterests,
  getinterestsProfile,
  avgPrice,
};

export default UserService;
