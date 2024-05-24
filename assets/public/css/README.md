
# MDBoostrap upgrade protocol

Before uploading a new version of MDboostrap some manual changes needs to be done. 

### Change the flags source file

You need to run the following command: 

```bash
sed -i 's#https://mdbootstrap.com/img/svg/flags.png#/assets/images/svg/flags.png#g' ./assets/public/css/mdb.min.css
```

The command above will force the framework to fetch the flags.png assets from our server instead of the mdboostrap.com server.
