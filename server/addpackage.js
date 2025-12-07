require('dotenv').config();

const mongoose = require('mongoose');
const Package = require('./Package');
const readline = require('node:readline');

const uri = process.env.AETHEIS_MONGODB_URI;

if (!uri) {
  console.error('FATAL ERROR: DB_URI is not defined.');
  process.exit(1); 
}

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
});

// Helper function to promisify readline.question
function askQuestion(query) {
    return new Promise(resolve => {
        rl.question(query, resolve);
    });
}

async function addPackage() {
    try {
        // 1. Connect to MongoDB
        await mongoose.connect(uri, { dbName: 'aetheis' });
        console.log('DB connected.');

        // 2. Collect input sequentially using await
        const name = await askQuestion('package name: ');
        // You might want to validate 'name' here before asking the next questions
        
        const installCmds = await askQuestion('package ic: ');
        const uninstallCmds = await askQuestion('package uc: ');

        let dependancies = [];
        while (true) {
            const dep = await askQuestion('add dependancy (leave blank to finish): ');
            if (dep.trim() === '') {
                break;
            }
            dependancies.push(dep.trim());
        }

        // 3. Close the readline interface after ALL input is collected
        rl.close();

        // 4. Create and Save the Package
        const newPkg = new Package({
            name: name,
            dependencies: dependancies,
            installCommands: installCmds,
            uninstallCommands: uninstallCmds
        });
        
        await newPkg.save();
        console.log('Package added successfully.');
        process.exit(0);

    } catch (e) {
        // Handle both DB connection and save errors
        console.error('Error adding package:', e);
        // Ensure readline interface is closed on error, if it's still open
        if (!rl.closed) {
            rl.close();
        }
        process.exit(1);
    }
}

addPackage();