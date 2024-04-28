import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import UserService from '../services/user.service';
const genderOptions = [
  { value: 2, label: "Male" },
  { value: 1, label: "Female" },
  { value: 0, label: "Other" },
  { value: 3, label: "Prefer Not to Say" }
];
const InterestsUserProfile = () => {
  const [profileData, setProfileData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const username = queryParams.get('username'); 
  useEffect(() => {
    const fetchProfileData = async () => {
      try {
        const data = await UserService.getinterestsProfile(username);
        setProfileData(data.data);
        setLoading(false);
      } catch (err) {
        console.error('Error fetching profile data:', err);
        setError('Failed to fetch profile data');
        setLoading(false);
      }
    };

    fetchProfileData();
  }, [username]); 

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>{error}</div>;
  }

  if (!profileData) {
    return <div>No profile data found for username: {username}</div>;
  }

  return (
    <div className="container">
      <header className="jumbotron">
        <h3>
          <strong>{profileData.username}</strong> Profile
        </h3>
      </header>
      <div>
        <p><strong>Username:</strong> {profileData.username}</p>
        <p><strong>First Name:</strong> {profileData.first_name}</p>
        <p><strong>Last Name:</strong> {profileData.last_name}</p>
        <p><strong>Date of Birth:</strong> {new Date(profileData.dob).toLocaleDateString('en-US', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit'
        })}</p>
        <p><strong>Gender:</strong> {genderOptions.find(option => option.value === profileData.gender)?.label}</p>
        <p><strong>Email:</strong> {profileData.email}</p>
        <p><strong>Phone:</strong> {profileData.phone}</p>
      </div>
    </div>
  );
};

export default InterestsUserProfile;
