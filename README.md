# Plugin for traefik: Convert Header

Replace some part of header by another value. Convert header to hex or int64 from uint64. Add prefix and postfix to header.

- 'fromHeader' - the header from which to take the value, required
- 'createHeader' - new header; remove the old value if it exists, required
- 'replaceValues' - replace oldValue by newValue, array, optional (First order)
- 'convertType' - for convert value from uint64 to hex use "uint64tohex", for convert from uint64toint65 use uint64toint64, only these values are available, optional (Second order)
- 'prefix' - add prefix to value, optional (Third order)
- 'postfix' - add postfix to value, optional (Fourth order)
  
Some operations can be skipped.

#### Yaml example:
```yaml
middlewares:
  convertHeader:
      plugin:
        convertheader:
          fromHeader: "X-OLD"
          createHeader: "X-NEW"
          replaceValues:
            - oldValue: "SerialNumber"
              newValue: ""
            - oldValue: "%"
              newValue: ""
          convertType: "uint64tohex"
          prefix: "SN="
          postfix: ";"
```

#### Labels example:
```
- traefik.http.middlewares.convertHeader.plugin.convertheader.fromHeader=X-OLD
- traefik.http.middlewares.convertHeader.plugin.convertheader.createHeader=X-NEW
- traefik.http.middlewares.convertHeader.plugin.convertheader.replaceValues[0].oldValue=SerialNumber
- traefik.http.middlewares.convertHeader.plugin.convertheader.replaceValues[0].newValue=
- traefik.http.middlewares.convertHeader.plugin.convertheader.replaceValues[1].oldValue=%
- traefik.http.middlewares.convertHeader.plugin.convertheader.replaceValues[1].newValue=
- traefik.http.middlewares.convertHeader.plugin.convertheader.convertType=uint64tohex
- traefik.http.middlewares.convertHeader.plugin.convertheader.prefix=SN
- traefik.http.middlewares.convertHeader.plugin.convertheader.postfix=;
```

Input: 'X-OLD: SerialNumber=9876543345678%'

Output: 'X-NEW: SN=8fb8fdb940e;'
