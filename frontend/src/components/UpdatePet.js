import React, { useState, useEffect } from 'react';
import UserService from '../services/user.service';

const UpdatePet = () => {
  const [petData, setPetData] = useState(null);
  const [message, setMessage] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        const data = await UserService.getpet();
        setPetData(data.data);
      } catch (error) {
        console.error('Error fetching pet data:', error);
        setMessage('Error fetching pet data');
      }
    };

    fetchData();
  }, []);

  const handleInputChange = (index, field, value) => {
    setPetData((prevPetData) => {
      const newData = [...prevPetData];
      newData[index][field] = value; 
      return newData;
    });
  };

  const handleUpdate = async (index) => {
    try {
      const petToUpdate = petData[index];
      await UserService.updatepet(
        petToUpdate.current_pet_name,
        petToUpdate.current_pet_type,
        petToUpdate.newPetName || '',
        petToUpdate.newPetType || '', 
        petToUpdate.newPetSize || ''
      );
      setMessage('Pet updated successfully');
    } catch (error) {
      console.error('Error updating pet:', error);
      if (error.response && error.response.data.error === 'DUPLICATE_KEY: pet already exists') {
        window.alert('DUPLICATE_KEY: pet already exists');
      } else {
        setMessage('Error updating pet');
      }
    }
  };
  

  return (
    <div className="container">
      {petData === null ? (
        <p>Loading...</p>
      ) : petData.length === 0 ? (
        <p>You do not have any pet.</p>
      ) : (
        petData.map((pet, index) => (
          <div key={index}>
            <h3>{`Pet ${index + 1}`}</h3>
            <p>
              <strong>Current Pet Name:</strong> {pet.current_pet_name}<br />
              <strong>Current Pet Type:</strong> {pet.current_pet_type}<br />
              <strong>Current Pet Size:</strong> {pet.current_pet_size}
            </p>
            <div>
              <label htmlFor={`newPetName${index}`}>New Pet Name:</label>
              <input
                type="text"
                id={`newPetName${index}`}
                value={pet.newPetName || ''}
                onChange={(e) => handleInputChange(index, 'newPetName', e.target.value)}
              />
            </div>
            <div>
              <label htmlFor={`newPetType${index}`}>New Pet Type:</label>
              <input
                type="text"
                id={`newPetType${index}`}
                value={pet.newPetType || ''}
                onChange={(e) => handleInputChange(index, 'newPetType', e.target.value)}
              />
            </div>
            <div>
              <label htmlFor={`newPetSize${index}`}>New Pet Size:</label>
              <input
                type="text"
                id={`newPetSize${index}`}
                value={pet.newPetSize || ''}
                onChange={(e) => handleInputChange(index, 'newPetSize', e.target.value)}
              />
            </div>
            <button onClick={() => handleUpdate(index)}>Update</button>
          </div>
        ))
      )}
      {message && <p>{message}</p>}
    </div>
  );
};

export default UpdatePet;
