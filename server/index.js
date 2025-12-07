require('dotenv').config();

const exp = require('express');
const mongoose = require('mongoose');
const Package = require('./Package');
const app = exp();

const uri = process.env.AETHEIS_MONGODB_URI;
const port = process.env.PORT || 3000;


if (!uri) {
  console.error('FATAL ERROR: DB_URI is not defined.');
  process.exit(1); 
}

// Connect to MongoDB
// Use 'uri' for the connection string and 'aetheis.packages' for the specific database
mongoose.connect(uri, {
  dbName: 'aetheis' 
    })
  .then(() => {
    console.log('DB connected.');
    
    // **SERVER STARTS HERE (ONLY ONCE)**
    app.listen(port, () => {
      console.log(`Server listening on port ${port}`);
    });
  })
  .catch((e) => {
    console.error('DB connection error:', e);
    process.exit(1);
  });

// Example route
app.get('/:package', async (req, res) => {
    const { package } = req.params;
    try {
        const pkgs = await Package.find({'name': package});
        if(pkgs.length === 0) {
            res.send("brew")
        } else {
            res.send("custom " + pkgs[0].dependencies.join(" "))
        }
    } catch (e) {
        res.status(500).json({ err: 'Failed to fetch packages' });
        console.error(e);
    }
        
});

app.get('/health', (req, res) => {
    res.send('Aetheis Package Server is running.');
});

app.get('/install/:package', async (req, res) => {
    const { package } = req.params;
    try {
        const pkg = await Package.findOne({'name': package});
        if(!pkg) {
            return res.status(404).json({ err: 'Package not found' });
        }
        res.send(pkg.installCommands)
    } catch (e) {
        res.status(500).json({ err: 'Failed to fetch package' });
    }
});

app.get('/uninstall/:package', async (req, res) => {
    const { package } = req.params;
    try {
        const pkg = await Package.findOne({'name': package});
        if(!pkg) {
            return res.status(404).json({ err: 'Package not found' });
        }
        res.send(pkg.uninstallCommands)
    } catch (e) {
        res.status(500).json({ err: 'Failed to fetch package' });
    }
});
// The duplicate app.listen() call has been removed from here.