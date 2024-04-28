import React, { useState } from 'react';
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import UserService from '../services/user.service';

const AvgPrice = () => {
  const [formData, setFormData] = useState({
    addr_zip_code: '',
    bedroom_num: '',
    bathroom_num: ''
  });
  const [avgRent, setAvgRent] = useState(null);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: name === 'bedroom_num' || name === 'bathroom_num' ? parseInt(value) : value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage('');
    setLoading(true);
  
    const { addr_zip_code, bedroom_num, bathroom_num } = formData;
  
    try {
      const response = await UserService.avgPrice(addr_zip_code, bedroom_num, bathroom_num); 
      setAvgRent(response.avg_rent); 
      setLoading(false);
    } catch (error) {
      if (error.response && error.response.data.error === "NO unit meets requirements") {
        setMessage('No unit meets requirements');
      } else {
        setMessage('Error retrieving average rent');
      }
      setLoading(false);
    }
  };

  return (
    <div className="container">
      <h2>Search Average Rent</h2>
      <Form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="addr_zip_code">Zip Code</label>
          <Input
            type="text"
            className="form-control"
            name="addr_zip_code"
            value={formData.addr_zip_code}
            onChange={handleChange}
          />
        </div>
        <div className="form-group">
          <label htmlFor="bedroom_num">Bedroom Count</label>
          <Input
            type="number"
            className="form-control"
            name="bedroom_num"
            value={formData.bedroom_num}
            onChange={handleChange}
          />
        </div>
        <div className="form-group">
          <label htmlFor="bathroom_num">Bathroom Count</label>
          <Input
            type="number"
            className="form-control"
            name="bathroom_num"
            onChange={handleChange}
          />
        </div>
        <div className="form-group">
          <Button className="btn btn-primary" disabled={loading}>Search</Button>
        </div>
      </Form>
      {message && <div className="alert alert-danger">{message}</div>} 
      {avgRent && <p>Average Rent: ${avgRent}</p>}
    </div>
  );
};

export default AvgPrice;
