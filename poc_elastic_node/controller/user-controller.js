const UserMongoRepository = require('../repository/user-mongo-repository');
const UserElasticRepository = require('../repository/user-elastic-repository');

class Controller {
  constructor(repository) {
    this.repository = repository;
  }

  async createUser(user) {
    return this.repository.createUser(user);
  }

  async getAllUsers(page = 1, size = 10) {
    return this.repository.getAllUsers(page, size);
  }

  async findUsersByName(name) {
    return this.repository.findUsersByName(name);
  }

  async findUsersByCity(city) {
    return this.repository.findUsersByCity(city);
  }

  async findUserByNRC(nrc) {
    return this.repository.findUserByNRC(nrc);
  }

  async createRandomUser(nrc) {
    const city = this.randomCity();
    const state = this.randomState(city);
    const randomUser = {
      name: this.randomName(),
      age: Math.floor(Math.random() * 100),
      nrc: nrc,
      date_of_registration: new Date(),
      address: {
        city: city,
        state: state,
        street: this.randomStreet(),
        number: Math.floor(Math.random() * 1000),
      },
    };

    await this.repository.createUser(randomUser);
    return randomUser;
  }

  async randomUpdate(nrc) {
    const user = await this.repository.findUserByNRC(nrc);
    if (!user) {
      throw new Error('User not found');
    }

    const city = this.randomCity();
    const state = this.randomState(city);
    user.name = this.randomName();
    user.age = Math.floor(Math.random() * 100);
    user.address.city = city;
    user.address.state = state;
    user.address.street = this.randomStreet();
    user.address.number = Math.floor(Math.random() * 1000);

    await this.repository.updateUser(user);
    return user;
  }

  randomName() {
    const names = ["João", "Maria", "Pedro", "Ana", "Carlos", "Fernanda", "Lucas", "Juliana", "Rafael", "Camila"];
    return names[Math.floor(Math.random() * names.length)];
  }

  randomCity() {
    const cities = ["São Paulo", "Rio de Janeiro", "Belo Horizonte", "Brasília", "Salvador", "Fortaleza", "Curitiba", "Recife", "Porto Alegre", "Manaus"];
    return cities[Math.floor(Math.random() * cities.length)];
  }

  randomState(city) {
    const cityStateMap = {
      "São Paulo": "SP",
      "Rio de Janeiro": "RJ",
      "Belo Horizonte": "MG",
      "Brasília": "DF",
      "Salvador": "BA",
      "Fortaleza": "CE",
      "Curitiba": "PR",
      "Recife": "PE",
      "Porto Alegre": "RS",
      "Manaus": "AM",
    };
    return cityStateMap[city];
  }

  randomStreet() {
    const streets = ["Rua das Flores", "Avenida Paulista", "Rua do Sol", "Avenida Atlântica", "Rua da Luz", "Avenida Brasil", "Rua da Paz", "Avenida das Américas", "Rua da Liberdade", "Avenida Central"];
    return streets[Math.floor(Math.random() * streets.length)];
  }
}

module.exports = Controller;
