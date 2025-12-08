const mongoose = require('mongoose');

// Define the Schema structure
const usrSch = new mongoose.Schema({
  name: {
    type: String,
    required: true, // This field must exist
    unique: true    // No two users can have the same name
  },
  password: String,
  token: String,
});

// Export the Schema wrapped in a Model
// 'User' is the Model name (singular, capitalized)
// Mongoose automatically looks for the plural, lowercase collection name: 'users'
module.exports = mongoose.model('User', usrSch);