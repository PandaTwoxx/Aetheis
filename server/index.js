require('dotenv').config();

const exp = require('express');
const mongoose = require('mongoose');
const Package = require('./Package');
const User = require('./User');
const app = exp();
const crypto = require('crypto');
const bcrypt = require('bcrypt');

const saltRounds = 10;

const uri = process.env.AETHEIS_MONGODB_URI;
const port = process.env.PORT || 3000;


if (!uri) {
    console.error('FATAL ERROR: DB_URI is not defined.');
    process.exit(1);
}


async function hashPassword(plainTextPassword) {
    try {
        const salt = await bcrypt.genSalt(saltRounds);
        const hash = await bcrypt.hash(plainTextPassword, salt);
        return hash;
    } catch (error) {
        console.error('Error hashing password:', error);
        throw error;
    }
}

async function verifyPassword(plainTextPassword, storedHashedPassword) {
    try {
        const isMatch = await bcrypt.compare(plainTextPassword, storedHashedPassword);
        return isMatch;
    } catch (error) {
        console.error('Error verifying password:', error);
        throw error;
    }
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
        const pkg = await Package.findOne({ 'name': package });
        console.log(pkg);
        if (!pkg) {
            res.send("brew")
        } else {
            if (!pkg.dependencies) {
                return res.send("custom")
            } else {
                res.send("custom " + pkg.dependencies.join(" "))
            }
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
        const pkg = await Package.findOne({ 'name': package });
        if (!pkg) {
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
        const pkg = await Package.findOne({ 'name': package });
        if (!pkg) {
            return res.status(404).json({ err: 'Package not found' });
        }
        res.send(pkg.uninstallCommands)
    } catch (e) {
        res.status(500).json({ err: 'Failed to fetch package' });
    }
});

app.get('/addUser/:user/:password', async (req, res) => {
    const { user, password } = req.params;
    try {
        const randomBytes = crypto.randomBytes(16);
        const hashedPassword = await hashPassword(password);
        const pkg = User.create({
            name: user,
            password: hashedPassword,
            token: randomBytes.toString('hex')
        });

        res.send(randomBytes.toString('hex'));
    } catch (e) {
        res.status(500).json({ err: 'Failed to create user' });
        console.error(e);
    }
});

app.get('/login/:user/:password', async (req, res) => {
    const { user, password } = req.params;
    try {
        const usr = User.findOne({ 'name': user });
        if (!usr) {
            return res.status(404).json({ err: 'User not found' });
        }
        const isMatch = await verifyPassword(password, usr.password);
        if (!isMatch) {
            return res.status(401).json({ err: 'Incorrect password' });
        }
        res.send(usr.token);
    } catch (e) {
        res.status(500).json({ err: 'Failed to login' });
        console.error(e);
    }
});

app.post('/addPackage/:token/:name/:installCmds/:uninstallCmds/:dependencies', async (req, res) => {
    const { token, name, installCmds, uninstallCmds, dependencies } = req.params;
    try {
        const usr = User.findOne({ 'token': token });
        if (!usr) {
            return res.status(404).json({ err: 'User not found' });
        }
        const pkg = Package.create({
            name: name,
            installCommands: installCmds,
            uninstallCommands: uninstallCmds,
            dependencies: dependencies.split(' '),
            owner: usr.name
        });
        res.send('Package added successfully.');
    } catch (e) {
        res.status(500).json({ err: 'Failed to add package' });
        console.error(e);
    }
});

app.post('/addPackage/:token/:name/:installCmds/:uninstallCmds/', async (req, res) => {
    const { token, name, installCmds, uninstallCmds } = req.params;
    try {
        const usr = User.findOne({ 'token': token });
        if (!usr) {
            return res.status(404).json({ err: 'User not found' });
        }
        const pkg = Package.create({
            name: name,
            installCommands: installCmds,
            uninstallCommands: uninstallCmds,
            dependencies: [],
            owner: usr.name
        });
        res.send('Package added successfully.');
    } catch (e) {
        res.status(500).json({ err: 'Failed to add package' });
        console.error(e);
    }
});

app.post('/updatePackage/:token/:name/:installCmds/:uninstallCmds/:dependencies', async (req, res) => {
    const { token, name, installCmds, uninstallCmds, dependencies } = req.params;
    try {
        const usr = User.findOne({ 'token': token });
        if (!usr) {
            return res.status(404).json({ err: 'User not found' });
        }
        const pkg = Package.findOne({ 'name': name, 'owner': usr.name });
        if (!pkg) {
            return res.status(404).json({ err: 'Package not found' });
        }
        pkg.installCommands = installCmds;
        pkg.uninstallCommands = uninstallCmds;
        pkg.dependencies = dependencies.split(' ');
        await pkg.save();
        res.send('Package updated successfully.');
    } catch (e) {
        res.status(500).json({ err: 'Failed to update package' });
        console.error(e);
    }
});

app.post('/deletePackage/:token/:name', async (req, res) => {
    const { token, name } = req.params;
    try {
        const usr = User.findOne({ 'token': token });
        if (!usr) {
            return res.status(404).json({ err: 'User not found' });
        }
        const pkg = Package.findOne({ 'name': name, 'owner': usr.name });
        if (!pkg) {
            return res.status(404).json({ err: 'Package not found' });
        }
        await pkg.deleteOne();
        res.send('Package deleted successfully.');
    } catch (e) {
        res.status(500).json({ err: 'Failed to delete package' });
        console.error(e);
    }
});
// The duplicate app.listen() call has been removed from here.