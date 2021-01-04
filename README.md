# Automated Toll Plaza

- This project contain a docker compose that will setup and run all the dependency for this application to run.
- Pre-requist to run this application requires an installation of docker.
- Once docker is install, change directory to docker and run the below commands
    - chmod +x start.sh
    - ./start.sh
- Once the script executes successfully, the application will be running on port 80.
- This project has the postman collection using which the APIs can be tested. It can be found in directory doc/postman
- To run the code coverage (or unit test cases) for the entire project, execute the below commands
    - chmod +x coverage.sh
    - ./coverage.sh
- This project has the swagger documentation for visual representation of the APIs request