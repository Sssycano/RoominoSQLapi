import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import UserService from '../services/user.service';

const Interests = () => {
  const [interestsData, setInterestsData] = useState(null);
  const [roommateCnt, setRoommateCnt] = useState(0);
  const [moveInDate, setMoveInDate] = useState('');
  const [complexUnitInfo, setComplexUnitInfo] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [createMessage, setCreateMessage] = useState('');

  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const unitRentId = parseInt(queryParams.get('unit_rent_id'));

  useEffect(() => {
    const fetchData = async () => {
      try {
        const interests = await UserService.getinterests(unitRentId);
        setInterestsData(interests.data);

        const complexInfo = await UserService.getcomplexunitinfo(unitRentId);
        setComplexUnitInfo(complexInfo.data);

        setLoading(false);
      } catch (err) {
        console.error('Error fetching data:', err);
        setError('Failed to fetch data');
        setLoading(false);
      }
    };

    fetchData();
  }, [unitRentId]);

  const handleCreateInterest = async (e) => {
    e.preventDefault();
    setCreateMessage('');

    try {
      const parsedRoommateCnt = parseInt(roommateCnt);
      const parsedMoveInDate = new Date(moveInDate);
      const formattedMoveInDate = parsedMoveInDate.toISOString().substring(0, 10);  // 保证格式为 'YYYY-MM-DD'
  
      await UserService.createinterests(unitRentId, parsedRoommateCnt, formattedMoveInDate);
      setCreateMessage('Interest created successfully!');
      
      setTimeout(() => {
        window.location.reload();
      }, 1500);
    } catch (err) {
      console.error('Error creating interest:', err);
      setCreateMessage('Failed to create interest.');
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>{error}</div>;
  }

  return (
    <div className="container">
      <h1>Building and Unit Information</h1>
      {complexUnitInfo && (
        <div>
          <h2>Building Information</h2>
          <p><strong>Company Name:</strong> {complexUnitInfo.company_name}</p>
          <p><strong>Building Name:</strong> {complexUnitInfo.building_name}</p>
          <p><strong>Address:</strong> {complexUnitInfo.address}</p>
          <p><strong>Year Built:</strong> {complexUnitInfo.year_built}</p>
          <p><strong>Available Units:</strong> {complexUnitInfo.available_units}</p>
          <p><strong>Building Amenities:</strong></p>
          <ul>
            {complexUnitInfo.amenities_building.map((amenity, index) => (
              <li key={index}>{amenity}</li>
            ))}
          </ul>

          <h2>Unit Information</h2>
          <p><strong>Unit Rent ID:</strong> {complexUnitInfo.unit_rent_id}</p> 
          <p><strong>Monthly Rent:</strong> {complexUnitInfo.monthly_rent}</p>
          <p><strong>Square Footage:</strong> {complexUnitInfo.square_footage}</p>
          <p><strong>Available Date for Move-In:</strong> {new Date(complexUnitInfo.available_date_for_move_in).toLocaleDateString()}</p>
          <p><strong>Bedrooms:</strong> {complexUnitInfo.bedroom_num}</p>
          <p><strong>Bathrooms:</strong> {complexUnitInfo.bathroom_num}</p>
          <p><strong>Living Rooms:</strong> {complexUnitInfo.living_room_num}</p>
          <p><strong>Unit Amenities:</strong></p>
          <ul>
            {complexUnitInfo.amenities_unit.map((amenity, index) => (
              <li key={index}>{amenity}</li>
            ))}
          </ul>
        </div>
      )}

      <h2>Current Interests</h2>
      {interestsData && interestsData.length > 0 ? (
        <ul>
          {interestsData.map((item, index) => (
            <li key={index}>
              <strong>Username:</strong> {item.username}<br />
              <strong>Unit Rent ID:</strong> {item.unit_rent_id}<br />
              <strong>Roommate Count:</strong> {item.roommate_cnt}<br />
              <strong>Move-in Date:</strong> {new Date(item.move_in_date).toLocaleDateString()}<br />
              <hr />
            </li>
          ))}
        </ul>
      ) : (
        <p>No interests found for unit {unitRentId}.</p>
      )}

      <h2>Create New Interest</h2>
      <form onSubmit={handleCreateInterest}>
        <div>
          <label htmlFor="roommateCnt">Roommate Count:</label>
          <input
            type="number"
            id="roommateCnt"
            value={roommateCnt}
            onChange={(e) => setRoommateCnt(e.target.value)}
          />
        </div>
        <div>
          <label htmlFor="moveInDate">Move-in Date:</label>
          <input
            type="date"
            id="moveInDate"
            value={moveInDate}
            onChange={(e) => setMoveInDate(e.target.value)}
          />
        </div>
        <button type="submit">Create Interest</button>
      </form>
      {createMessage && <p>{createMessage}</p>}
    </div>
  );
};

export default Interests;
