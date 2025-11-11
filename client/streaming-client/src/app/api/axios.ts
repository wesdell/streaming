import axios from 'axios';

const apiURL: string = import.meta.env.VITE_API_URL;

export default axios.create({
  baseURL: apiURL,
  headers: {
    'Content-Type': 'application/json'
  }
})
