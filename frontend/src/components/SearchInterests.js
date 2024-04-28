import React, { useState } from 'react';
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import UserService from '../services/user.service';
import { Link } from 'react-router-dom'; 
const SearchInterests = () => {
  const [formData, setFormData] = useState({
    unit_rent_id: '',
    roommate_cnt: '',
    move_in_date: ''
  });
  const [searchResult, setSearchResult] = useState([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');

  const handleChange = (e) => {
    const { name, value } = e.target;
    if (name === 'unit_rent_id' || name === 'roommate_cnt') {
      setFormData({ ...formData, [name]: value ? parseInt(value) : '' });
    } else {
      setFormData({ ...formData, [name]: value });
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    setMessage('');
    setLoading(true);
    const filteredFormData = {};
    Object.keys(formData).forEach((key) => {
      if (formData[key] !== '') {
        filteredFormData[key] = formData[key];
      }
    });

    UserService.searchInterests(
      filteredFormData.unit_rent_id,
      filteredFormData.roommate_cnt,
      filteredFormData.move_in_date
    )
      .then((data) => {
        if (data.data === null) {
          setSearchResult([]);
          setMessage('No Result');
        } else {
          setSearchResult(data.data);
          setMessage('');
        }

        setLoading(false);
      })
      .catch((error) => {
        setMessage('Error retrieving search results');
        setLoading(false);
      });
  };

  return (
    <div className="container">
      <div className="row">
        <div className="col-md-6 offset-md-3">
          <h2>Search Interests</h2>
          <Form onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="unit_rent_id">Unit Rent ID</label>
              <Input
                type="number"
                className="form-control"
                name="unit_rent_id"
                value={formData.unit_rent_id}
                onChange={handleChange}
              />
            </div>
            <div className="form-group">
              <label htmlFor="roommate_cnt">Roommate Count</label>
              <Input
                type="number"
                className="form-control"
                name="roommate_cnt"
                value={formData.roommate_cnt}
                onChange={handleChange}
              />
            </div>
            <div className="form-group">
              <label htmlFor="move_in_date">Move-in Date</label>
              <Input
                type="date"
                className="form-control"
                name="move_in_date"
                value={formData.move_in_date}
                onChange={handleChange}
              />
            </div>
            <div className="form-group">
              <Button className="btn btn-primary" disabled={loading}>Search</Button>
            </div>
          </Form>
          {message && <div className="alert alert-danger">{message}</div>}
          <hr />
          <div>
            {searchResult.length > 0 ? (
              <div>
                <h3>Search Results</h3>
                <ul>
                  {searchResult.map((result) => (
                    <li key={result.username}>
                      <Link to={`/profile/interestprofile?username=${result.username}`}> 
                        {result.username}
                      </Link>
                      <strong>Unit Rent ID:</strong> {result.unit_rent_id}<br />
                      <strong>Roommate Count:</strong> {result.roommate_cnt}<br />
                      <strong>Move-in Date:</strong> {new Date(result.move_in_date).toLocaleDateString('en-US', {
                        year: 'numeric',
                        month: '2-digit',
                        day: '2-digit'
                      })}<br />
                      <hr />
                    </li>
                  ))}
                </ul>
              </div>
            ) : (
              <div>
                <h3>Search Results</h3>
                <p>No Result</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default SearchInterests;
