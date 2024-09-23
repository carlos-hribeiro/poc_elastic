const { MongoClient } = require('mongodb');

class UserMongoRepository {
  constructor(client) {
    this.client = client;
  }

  async createUser(user) {
    user.date_of_registration = new Date();
    const result = await this.client.db('test').collection('users').insertOne(user);
    return result;
  }

  async updateUser(user) {
    const result = await this.client.db('test').collection('users').updateOne({ id: user.id }, { $set: user });
    return result;
  }

  async getAllUsers(page = 1, size = 10) {
    const skip = (page - 1) * size;
    const users = await this.client.db('test').collection('users').find().sort({ date_of_registration: -1 }).skip(skip).limit(size).toArray();
    return users;
  }

  async findUsersByName(name) {
    const users = await this.client.db('test').collection('users').find({ name: { $regex: name, $options: 'i' } }).toArray();
    return users;
  }

  async findUsersByCity(city) {
    const users = await this.client.db('test').collection('users').find({ 'address.city': { $regex: city, $options: 'i' } }).toArray();
    return users;
  }

  async findUserByNRC(nrc) {
    const user = await this.client.db('test').collection('users').findOne({ nrc });
    return user;
  }
}

module.exports = UserMongoRepository;
