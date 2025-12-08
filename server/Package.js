const mongoose = require('mongoose');

// Define the Schema structure
const pkgSch = new mongoose.Schema({
  name: {
    type: String,
    required: true, // This field must exist
    unique: true    // No two packages can have the same name
  },
  installCommands: String,
    uninstallCommands: String,
    owner: String
});

// Export the Schema wrapped in a Model
// 'Package' is the Model name (singular, capitalized)
// Mongoose automatically looks for the plural, lowercase collection name: 'packages'
module.exports = mongoose.model('Package', pkgSch);