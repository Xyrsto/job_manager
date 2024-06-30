# Job Manager - First version is up!
![imagem](https://github.com/Xyrsto/job_manager/assets/73367973/98d4c4cb-1430-459b-84ed-1449d882d32a)
## How to use?
Since authentication is not implemented, you have to use your own MongoDB database to use this version.
- 1st step: Create a .env file in the folder where you have the JobManager.exe.
- 2nd step: Setup the MongoDB database:
    - Start you cluster
    - Create a user to manage the database.
    - Then, on the left side click on "Database" and click "Connect" to copy the URI.
    - Click on "Browse Collections" and the "+ Create Database". Name it whatever you want.
    - On the database you have created, create a collection. Again, set whatever name you want.
 - 3d step: On the .env file you have created, create 3 environment variables: :
    - "MONGO_DB_URI": copy and paste your MongoDB URI here.
    - "MONGO_DB_DATABASE": here you should put the name of the database you have created.
    - "MONGO_DB_COLLECTION": here you should put the name of the collection you have created.
Finally, your app should launch without any problems!

## What commands are available?
Since i haven't implemented any commands to help you use this app (I am planning on doing it in the future), here is the list of commands supported by it:
- `jm -a`: This one is used to create an entry on your job table. 
    - `-cn`: Used to set the company name.
    - `-r`: Used to set the rating of the company
    - `-n`: Used to set a short note about the job offer.
    - `-ha`: Used to set if you have received an answer about the job offer, like a job interview confirmation.
        - This one can either be `true` or `false`. The app will filter your table to show your confirmed job interviews first (the ones set to true)
    - **Example**:
        - `jm -a -cn Company Name -r 5/5 -n Very cool offer! -ha false`
- `jm -ls`: This command will list all the job offers you have registered in the database, in a pretty little table.
- `jm -u`: Using this command, you will need to provide the name of the company you want to update. What it will, is update the "Has Answered" field to `true`

## What's planned for the next version
- Update more offer fields
- Help command
- Delete an offer



