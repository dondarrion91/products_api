#!/usr/bin/env bash

set -e

echo "📦 Generando datos de prueba (20 categorías, 20 sellers, 100 imágenes, 100 productos)..."

#############################
# CATEGORÍAS
#############################

CATEGORY_NAMES=(
  "Electrónica"
  "Hogar y Deco"
  "Cocina"
  "Gaming"
  "Oficina"
  "Deportes"
  "Jardín"
  "Belleza"
  "Salud"
  "Automotor"
  "Mascotas"
  "Bebés"
  "Librería"
  "Audio y Música"
  "Fotografía"
  "Climatización"
  "Iluminación"
  "Herramientas"
  "Moda"
  "Viajes y Outdoor"
)

CATEGORY_DESCRIPTIONS=(
  "Dispositivos electrónicos, celulares y computación."
  "Productos para el hogar, decoración y muebles."
  "Accesorios y electrodomésticos para la cocina."
  "Periféricos y accesorios gamer."
  "Sillas, escritorios y accesorios de oficina."
  "Indumentaria y accesorios deportivos."
  "Herramientas y productos para el jardín."
  "Productos de cuidado personal y belleza."
  "Artículos relacionados con la salud y bienestar."
  "Accesorios y repuestos para autos."
  "Accesorios y alimento para mascotas."
  "Productos para bebés y primera infancia."
  "Libros, cuadernos y artículos escolares."
  "Parlantes, auriculares y equipos de audio."
  "Cámaras y accesorios de fotografía."
  "Aires acondicionados, ventiladores y calefacción."
  "Lámparas y sistemas de iluminación."
  "Herramientas manuales y eléctricas."
  "Ropa, calzado y accesorios."
  "Mochilas, valijas y equipo outdoor."
)

echo "🗂 Creando Category.json"
{
  echo "["
  for ((i=0; i<20; i++)); do
    id=$((i+1))
    name="${CATEGORY_NAMES[$i]}"
    desc="${CATEGORY_DESCRIPTIONS[$i]}"
    if (( i == 19 )); then comma=""; else comma=","; fi
    cat <<EOF
  {
    "id": "cat-$id",
    "name": "$name",
    "description": "$desc"
  }$comma
EOF
  done
  echo "]"
} > Category.json

#############################
# SELLERS
#############################

SELLER_NAMES=(
  "TechWorld Store"
  "HomeCenter"
  "CasaMarket"
  "GamerZone"
  "ElectroCity"
  "OfficePlus"
  "GreenGarden"
  "BeautyCare Shop"
  "HealthPlus"
  "AutoParts Max"
  "PetLovers"
  "BabySmile"
  "Book&Paper"
  "SoundWave"
  "PhotoPro"
  "ClimaCool"
  "LightHouse"
  "ToolMaster"
  "UrbanStyle"
  "OutdoorLife"
)

echo "🧑‍💼 Creando Seller.json"
{
  echo "["
  for ((i=0; i<20; i++)); do
    id=$((i+1))
    name="${SELLER_NAMES[$i]}"
    rating=$(( (RANDOM % 2) + 4 )) # 4 o 5
    if (( i == 19 )); then comma=""; else comma=","; fi
    cat <<EOF
  {
    "id": "seller-$id",
    "name": "$name",
    "rating": $rating
  }$comma
EOF
  done
  echo "]"
} > Seller.json

#############################
# IMÁGENES
#############################

echo "🖼 Creando Image.json"
{
  echo "["
  for ((i=0; i<100; i++)); do
    id=$((i+1))
    url="https://picsum.photos/seed/img$id/600/600"
    if (( i == 99 )); then comma=""; else comma=","; fi
    cat <<EOF
  {
    "id": "img-$id",
    "url": "$url"
  }$comma
EOF
  done
  echo "]"
} > Image.json

#############################
# PRODUCTOS
#############################

PRODUCT_NAMES=(
  "Auriculares Bluetooth"
  "Mouse Gamer RGB"
  "Teclado Mecánico"
  "Monitor 24 pulgadas"
  "Silla Gamer"
  "Notebook Ultrabook"
  "Smartwatch Deportivo"
  "Parlante Bluetooth"
  "Lámpara LED de escritorio"
  "Sartén antiadherente"
  "Set de cuchillos de cocina"
  "Aspiradora inalámbrica"
  "Silla ergonómica"
  "Mochila para notebook"
  "Cámara web HD"
  "Micrófono USB"
  "Disco externo 1TB"
  "Memoria USB 64GB"
  "Router WiFi"
  "Cafetera eléctrica"
)

PRODUCT_DESCRIPTIONS=(
  "Auriculares inalámbricos con cancelación de ruido y estuche de carga."
  "Mouse gamer con sensor óptico de alta precisión y luces RGB."
  "Teclado mecánico con switches táctiles y retroiluminación."
  "Monitor LED Full HD de 24 pulgadas, ideal para oficina y gaming casual."
  "Silla gamer con soporte lumbar y reclinación ajustable."
  "Notebook ultraliviana ideal para trabajo y estudio."
  "Reloj inteligente con monitoreo de actividad y notificaciones."
  "Parlante portátil Bluetooth con batería de larga duración."
  "Lámpara LED de escritorio con brillo regulable."
  "Sartén con recubrimiento antiadherente de alta resistencia."
  "Juego de cuchillos de cocina de acero inoxidable."
  "Aspiradora inalámbrica para limpieza rápida del hogar."
  "Silla ergonómica para largas jornadas de trabajo."
  "Mochila con compartimento acolchado para notebook."
  "Cámara web HD para videollamadas y streaming."
  "Micrófono USB para podcasting y videollamadas."
  "Disco rígido externo de 1TB para backup."
  "Pendrive USB 3.0 de 64GB."
  "Router WiFi de doble banda."
  "Cafetera eléctrica para café filtrado."
)

DETAIL_COLOR=( "Negro" "Blanco" "Rojo" "Azul" "Gris" "Verde" "Plateado" "Dorado" )
DETAIL_GARANTIA=( "6 meses" "12 meses" "18 meses" "24 meses" )

echo "📦 Creando Product.json"
{
  echo "["
  for ((i=0; i<100; i++)); do
    id=$((i+1))

    # Nombre y descripción basados en arrays
    idx=$(( i % ${#PRODUCT_NAMES[@]} ))
    name="${PRODUCT_NAMES[$idx]}"
    desc="${PRODUCT_DESCRIPTIONS[$idx]}"

    # Category y seller cíclicos
    catId=$(( (i % 20) + 1 ))
    sellerId=$(( (i % 20) + 1 ))

    # Rate 3–5
    rate=$(( (RANDOM % 3) + 3 ))

    # Precio 19–999
    price=$(( (RANDOM % 900) + 100 ))

    # Descuento 0, 5, 10, 15, 20
    discounts=(0 5 10 15 20)
    discIdx=$(( RANDOM % ${#discounts[@]} ))
    discount=${discounts[$discIdx]}

    # Cuotas 1, 3, 6, 12
    installments_list=(1 3 6 12)
    instIdx=$(( RANDOM % ${#installments_list[@]} ))
    installments=${installments_list[$instIdx]}

    # Stock 0–50
    stock=$(( RANDOM % 51 ))

    # Ventas 0–500
    sales=$(( RANDOM % 501 ))

    # Details
    color="${DETAIL_COLOR[$((RANDOM % ${#DETAIL_COLOR[@]}))]}"
    garantia="${DETAIL_GARANTIA[$((RANDOM % ${#DETAIL_GARANTIA[@]}))]}"

    # Imágenes: 1–3 imágenes por producto
    img1=$(( (i % 100) + 1 ))
    img2=$(( (i*7 % 100) + 1 ))
    img3=$(( (i*13 % 100) + 1 ))
    num_imgs=$(( (RANDOM % 3) + 1 ))

    if (( num_imgs == 1 )); then
      images_json="\"img-$img1\""
    elif (( num_imgs == 2 )); then
      images_json="\"img-$img1\", \"img-$img2\""
    else
      images_json="\"img-$img1\", \"img-$img2\", \"img-$img3\""
    fi

    # Characteristics simple
    char_name="Características del producto"

    if (( i == 99 )); then comma=""; else comma=","; fi

    cat <<EOF
  {
    "id": "prod-$id",
    "name": "$name",
    "rate": $rate,
    "price": $price,
    "discount": $discount,
    "installments": $installments,
    "stock": $stock,
    "details": [
      { "name": "Color", "description": "$color" },
      { "name": "Garantía", "description": "$garantia" }
    ],
    "images": [
      $images_json
    ],
    "sales_number": $sales,
    "description": "$desc",
    "categoryId": "cat-$catId",
    "sellerId": "seller-$sellerId",
    "characteristics": {
      "name": "$char_name",
      "details": null
    }
  }$comma
EOF
  done
  echo "]"
} > Product.json

echo "✅ Listo. Se generaron:"
echo "   - Category.json (20 categorías)"
echo "   - Seller.json   (20 sellers)"
echo "   - Image.json    (100 imágenes)"
echo "   - Product.json  (100 productos)"
