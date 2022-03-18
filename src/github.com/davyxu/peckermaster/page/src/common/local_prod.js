
let Host=''

if (process.env.NODE_ENV ==='development'){
  Host='http://localhost:9096'
}

export default{
  Host
}
