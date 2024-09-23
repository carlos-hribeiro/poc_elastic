const express = require('express');
const bodyParser = require('body-parser');
const MongoClient = require('mongodb').MongoClient;
const { Client: ElasticClient } = require('@elastic/elasticsearch');
const Controller = require('./controller/user-controller');
const UserMongoRepository = require('./repository/user-mongo-repository');
const UserElasticRepository = require('./repository/user-elastic-repository');

async function main(){

  const app = express();
  app.use(bodyParser.json());

  const backend = process.argv[2] || 'elastic';
  let repository;

  console.log(`Backend: ${backend}`);

  if (backend === 'elastic') {
    const elasticClient = new ElasticClient({ node: 'https://localhost:9200', auth: { username: 'elastic', password: 'G+Rehd00aiBg8KpPKHNf' }, tls: { rejectUnauthorized: false } });
    repository = new UserElasticRepository(elasticClient);
  } else if (backend === 'mongo') {
    try{
    let client = await MongoClient.connect('mongodb://localhost:27017', { useNewUrlParser: true, useUnifiedTopology: true });
   
    repository = new UserMongoRepository(client);
    
    }catch(err){
      console.error('Error creating the client:', err);
      process.exit(1);
    }

  } else {
    console.error('Invalid backend:', backend);
    process.exit(1);
  }

  const controller = new Controller(repository);

  app.post('/users', async (req, res) => {
    try {
      const user = req.body;
      await controller.createUser(user);
      res.status(201).send('Client created successfully');
    } catch (err) {
      console.error(err);
      res.status(500).send('Error saving user to database');
    }
  });

  app.get('/users/all', async (req, res) => {
    try {
      const page = parseInt(req.query.page) || 1;
      const size = parseInt(req.query.size) || 10;
      const users = await controller.getAllUsers(page, size);
      res.json(users);
    } catch (err) {
      console.error(err);
      res.status(500).send('Error fetching users from database');
    }
  });

  app.get('/users/findByName', async (req, res) => {
    try {
      const name = req.query.name;
      if (!name) {
        res.status(400).send('Name parameter is required');
        return;
      }
      const users = await controller.findUsersByName(name);
      res.json(users);
    } catch (err) {
      console.error(err);
      res.status(500).send('Error fetching users from database');
    }
  });

  app.get('/users/findByCity', async (req, res) => {
    try {
      const city = req.query.city;
      if (!city) {
        res.status(400).send('City parameter is required');
        return;
      }
      const users = await controller.findUsersByCity(city);
      res.json(users);
    } catch (err) {
      console.error(err);
      res.status(500).send('Error fetching users from database');
    }
  });

  app.get('/users/findByNRC', async (req, res) => {
    try {
      const nrc = parseInt(req.query.nrc);
      if (isNaN(nrc)) {
        res.status(400).send('Invalid NRC parameter');
        return;
      }
      const user = await controller.findUserByNRC(nrc);
      if (!user) {
        res.status(404).send('User not found');
        return;
      }
      res.json(user);
    } catch (err) {
      console.error(err);
      res.status(500).send('Error fetching user from database');
    }
  });

  app.post('/users/random', async (req, res) => {
    try {
      const nrc = req.body.nrc;
      const user = await controller.createRandomUser(nrc);
      res.json(user);
    } catch (err) {
      console.error(err);
      res.status(500).send('Error creating random user');
    }
  });

  app.post('/users/random-update', async (req, res) => {
    try {
      const nrc = req.body.nrc;
      const user = await controller.randomUpdate(nrc);
      res.json(user);
    } catch (err) {
      console.error(err);
      res.status(500).send('Error updating user');
    }
  });

  const PORT = 8080;
  app.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
  });


}


main();