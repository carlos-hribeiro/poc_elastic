const { Client } = require('@elastic/elasticsearch');

class UserElasticRepository {
  constructor(client) {
    this.client = client;
  }

  async createUser(user) {
    user.date_of_registration = new Date();
    const body = await this.client.index({
      index: 'users',
      body: user,
    });
    return body;
  }

  async updateUser(user) {
    const body = await this.client.update({
      index: 'users',
      id: user.id,
      body: { doc: user },
    });
    return body;
  }

  async getAllUsers(page = 1, size = 10) {
    const body = await this.client.search({
      index: 'users',
      sort: 'date_of_registration:desc',
      from: (page - 1) * size,
      size: size,
    });
    return body.hits.hits.map(hit => ({ ...hit._source, id: hit._id }));
  }

  async findUsersByName(name) {
    const body = await this.client.search({
      index: 'users',
      body: {
        query: {
          prefix: { name: name },
        },
      },
    });
    return body.hits.hits.map(hit => ({ ...hit._source, id: hit._id }));
  }

  async findUsersByCity(city) {
    const body = await this.client.search({
      index: 'users',
      body: {
        query: {
          prefix: { 'address.city': city },
        },
      },
    });
    return body.hits.hits.map(hit => ({ ...hit._source, id: hit._id }));
  }

  async findUserByNRC(nrc) {
    const body = await this.client.search({
      index: 'users',
      body: {
        query: {
          term: { nrc: nrc },
        },
      },
    });
    if (body.hits.total.value === 0) {
      return null;
    }
    const user = body.hits.hits[0];
    return { ...user._source, id: user._id };
  }
}

module.exports = UserElasticRepository;
